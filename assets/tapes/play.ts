// Play the real app with normal keys to capture grown snakes and Tetris line clears.
// The tmux server is private to this recording; the app and saved preferences are unchanged.
import { mkdtemp, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";

type Cell = { x: number; y: number };
const socket = `arcade-vhs-${process.pid}`;
const game = process.argv[2] ?? "tetris";
if (game !== "tetris" && game !== "snake") throw new Error("Choose tetris or snake");
const dataPath = await mkdtemp(join(tmpdir(), "arcade-vhs-"));
let stopping = false;

async function tmux(...args: string[]) {
	const child = Bun.spawn(["tmux", "-L", socket, ...args], {
		stdout: "pipe", stderr: "pipe",
	});
	const [output, error, code] = await Promise.all([
		new Response(child.stdout).text(), new Response(child.stderr).text(), child.exited,
	]);
	if (code !== 0) throw new Error(error.trim() || `tmux exited with ${code}`);
	return output;
}

async function screen() {
	return (await tmux("capture-pane", "-p", "-t", "demo"));
}

async function waitFor(text: string) {
	const deadline = Date.now() + 20_000;
	while (!stopping && Date.now() < deadline) {
		if ((await screen()).includes(text)) return;
		await Bun.sleep(80);
	}
	throw new Error(`The app did not show ${text}`);
}

function readBoard(view: string) {
	const rows = view.split("\n");
	const top = rows.findIndex(row => row.includes("┌────────────────────┐"));
	if (top < 0) throw new Error("Tetris board is not visible");
	const left = rows[top].indexOf("┌") + 1;
	const cells: Cell[] = [];
	for (let y = 0; y < 20; y++) {
		for (let x = 0; x < 10; x++) {
			if (rows[top + y + 1]?.[left + 2 * x] === "█") cells.push({ x, y });
		}
	}
	const active = cells.slice(0, 4);
	if (active.length !== 4) throw new Error("The current piece is not visible");
	const board = Array.from({ length: 20 }, () => Array<boolean>(10).fill(false));
	for (const { x, y } of cells.slice(4)) board[y][x] = true;
	return { active, board };
}

function normalize(cells: Cell[]) {
	const left = Math.min(...cells.map(cell => cell.x));
	const top = Math.min(...cells.map(cell => cell.y));
	return cells.map(cell => ({ x: cell.x - left, y: cell.y - top }));
}

// Prefer complete rows, low stacks, and empty cells that remain reachable from above.
function chooseMove(board: boolean[][], active: Cell[]) {
	let shape = normalize(active);
	let best = { score: -Infinity, rotation: 0, x: 0 };
	for (let rotation = 0; rotation < 4; rotation++) {
		const width = 1 + Math.max(...shape.map(cell => cell.x));
		for (let x = 0; x <= 10 - width; x++) {
			const fits = (y: number) => shape.every(cell => y + cell.y < 20 && !board[y + cell.y][x + cell.x]);
			if (!fits(0)) continue;
			let y = 0;
			while (fits(y + 1)) y++;
			const result = board.map(row => [...row]);
			for (const cell of shape) result[y + cell.y][x + cell.x] = true;
			const rows = result.filter(row => !row.every(Boolean));
			const cleared = 20 - rows.length;
			while (rows.length < 20) rows.unshift(Array<boolean>(10).fill(false));
			let holes = 0;
			const heights = Array.from({ length: 10 }, (_, column) => {
				const top = rows.findIndex(row => row[column]);
				if (top < 0) return 0;
				for (let row = top; row < 20; row++) if (!rows[row][column]) holes++;
				return 20 - top;
			});
			const height = heights.reduce((sum, value) => sum + value, 0);
			const bumps = heights.slice(1).reduce((sum, value, i) => sum + Math.abs(value - heights[i]), 0);
			const score = 8 * cleared - 0.5 * height - 4 * holes - 0.3 * bumps;
			if (score > best.score) best = { score, rotation, x };
		}
		const height = 1 + Math.max(...shape.map(cell => cell.y));
		shape = normalize(shape.map(cell => ({ x: height - 1 - cell.y, y: cell.x })));
	}
	return best;
}

async function playTetris() {
	await waitFor("Theme:");
	await Bun.sleep(1500);
	await tmux("send-keys", "-t", "demo", "Down");
	await Bun.sleep(250);
	await tmux("send-keys", "-t", "demo", "Down");
	await Bun.sleep(450);
	await tmux("send-keys", "-t", "demo", "Enter");
	await waitFor("T E T R I S");
	await Bun.sleep(700);
	for (let piece = 0; piece < 30 && !stopping; piece++) {
		const view = await screen();
		if (view.includes("GAME OVER")) throw new Error("The recording reached game over");
		const { active, board } = readBoard(view);
		const move = chooseMove(board, active);
		for (let turn = 0; turn < move.rotation; turn++) {
			await tmux("send-keys", "-t", "demo", "Up");
			await Bun.sleep(100);
		}
		await Bun.sleep(80);
		const rotated = readBoard(await screen()).active;
		const distance = move.x - Math.min(...rotated.map(cell => cell.x));
		for (let step = 0; step < Math.abs(distance); step++) {
			await tmux("send-keys", "-t", "demo", distance < 0 ? "Left" : "Right");
			await Bun.sleep(75);
		}
		await Bun.sleep(200);
		await tmux("send-keys", "-t", "demo", "Space");
		await Bun.sleep(450);
	}
}

async function playSnake() {
	await waitFor("Theme:");
	await tmux("send-keys", "-t", "demo", "Down", "Enter");
	await waitFor("S N A K E");
	let previous = "";
	const moves = [
		{ x: 1, y: 0, key: "Right" }, { x: 0, y: 1, key: "Down" },
		{ x: -1, y: 0, key: "Left" }, { x: 0, y: -1, key: "Up" },
	];
	const deadline = Date.now() + 60_000;
	while (!stopping && Date.now() < deadline) {
		const view = await screen();
		if (view.includes("GAME OVER")) throw new Error("Snake reached game over");
		const rows = view.split("\n");
		const top = rows.findIndex(row => row.includes("┌" + "─".repeat(56) + "┐"));
		if (top < 0) throw new Error("Snake board is not visible");
		const left = rows[top].indexOf("┌") + 1;
		const blocked = new Set<string>();
		let head: Cell | undefined;
		let food: Cell | undefined;
		for (let y = 0; y < 16; y++) {
			for (let x = 0; x < 28; x++) {
				const glyph = rows[top + y + 1]?.[left + 2 * x];
				if (glyph === "█" || glyph === "▓") blocked.add(`${x},${y}`);
				if (glyph === "▓") head = { x, y };
				if (glyph === "◆") food = { x, y };
			}
		}
		if (!head || !food) throw new Error("Snake or food is not visible");
		const body = [...blocked].map(cell => cell.split(",").map(Number));
		const horizontalSpan = Math.max(...body.map(cell => cell[0])) - Math.min(...body.map(cell => cell[0]));
		const verticalSpan = Math.max(...body.map(cell => cell[1])) - Math.min(...body.map(cell => cell[1]));
		if (Number(view.match(/SCORE (\d+)/)?.[1]) >= 8 && horizontalSpan >= 3 && verticalSpan >= 3) {
			await tmux("send-keys", "-t", "demo", "Space");
			return;
		}
		const position = `${head.x},${head.y}`;
		if (position === previous) { await Bun.sleep(15); continue; }
		previous = position;
		const queue = [{ ...head, first: "" }];
		const visited = new Set([position]);
		let next = "";
		for (let i = 0; i < queue.length && !next; i++) {
			const cell = queue[i];
			if (cell.x === food.x && cell.y === food.y) { next = cell.first; break; }
			for (const move of moves) {
				const x = cell.x + move.x, y = cell.y + move.y;
				const key = `${x},${y}`;
				if (x < 0 || x >= 28 || y < 0 || y >= 16 || blocked.has(key) || visited.has(key)) continue;
				visited.add(key);
				queue.push({ x, y, first: cell.first || move.key });
			}
		}
		if (!next) throw new Error("Food is not reachable in this recording");
		await tmux("send-keys", "-t", "demo", next);
		await Bun.sleep(15);
	}
	if (!stopping) throw new Error("Snake did not reach the target score");
}

try {
	await tmux("new-session", "-d", "-s", "demo", "-x", "80", "-y", "24", "-e", `XDG_DATA_HOME=${dataPath}`, "go run .");
	await tmux("set-option", "-g", "status", "off");
	const attached = Bun.spawn(["tmux", "-L", socket, "attach-session", "-t", "demo"], {
		stdin: "inherit", stdout: "inherit", stderr: "inherit",
	});
	const playing = (game === "snake" ? playSnake() : playTetris()).catch(error => {
		if (!stopping) {
			console.error(error);
			process.exitCode = 1;
			attached.kill();
		}
	});
	await attached.exited;
	stopping = true;
	await playing;
} finally {
	stopping = true;
	await tmux("kill-server").catch(() => {});
	await rm(dataPath, { recursive: true, force: true });
}
