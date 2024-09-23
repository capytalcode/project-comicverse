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
	pages: {
		title: string;
		src: string;
		background: string;
	}[];
};

export type { Project };
export default db;
