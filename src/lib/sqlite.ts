import sqlite3 from 'sqlite3';
import { open } from 'sqlite';

const db = await open({
	filename: 'data.db',
	driver: sqlite3.cached.Database
});

await db.exec(`
CREATE TABLE IF NOT EXISTS projects (
	ID text NOT NULL,
	Name text NOT NULL,
	PRIMARY KEY(ID)
)
`);

type Project = {
	id: string;
	title: string;
	pages: Page[];
};

type Page = {
	title: string;
	src: string;
	background: string;
	iteraction: Iteraction[];
};

type Iteraction = {
	x: number;
	y: number;
	link: string;
};

export type { Project, Iteraction, Page };
export default db;
