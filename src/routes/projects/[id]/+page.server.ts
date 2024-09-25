import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import stream from 'node:stream/promises';
import { db, s3, type Project } from '$lib';
import { AWS_S3_DEFAULT_BUCKET } from '$env/static/private';
import { extname } from 'node:path';

export const prerender = false;
export const ssr = false;

export const load = (async ({ params }) => {
	const res = await db.get<{ id: string; name: string }>(
		'SELECT ID, Name FROM projects WHERE ID = ?',
		params.id
	);

	if (res === undefined) {
		return fail(404, { reason: 'Failed to find project into database' });
	}

	const project = await s3.getObject(AWS_S3_DEFAULT_BUCKET, `${params.id}/project.json`);
	project.on('error', (err: any) => {
		console.log(err);
	});
	let p: string = '';
	project.on('data', (chunk: any) => {
		p += chunk;
	});
	await stream.finished(project);
	let proj = JSON.parse(p) as Project;

	return { project: proj };
}) as PageServerLoad;

export const actions = {
	delete: async ({ params }) => {
		const res = await db.run('DELETE FROM projects WHERE ID = ?', params.id);

		await s3.removeObject(AWS_S3_DEFAULT_BUCKET, `${params.id}/project.json`);

		if (res === undefined) {
			return fail(500, { reason: 'Failed to delete project' });
		}

		redirect(303, '/');
	},
	addpage: async ({ request, params }) => {
		const form = await request.formData();
		const file = form?.get('file') as File;
		const title = form?.get('title') as string;
		const color = form?.get('color') as string;

		console.log(file);

		const project = await s3.getObject(AWS_S3_DEFAULT_BUCKET, `${params.id}/project.json`);
		project.on('error', (err: any) => {
			console.log(err);
		});
		let p: string = '';
		project.on('data', (chunk: any) => {
			p += chunk;
		});
		await stream.finished(project);
		let proj = JSON.parse(p) as Project;

		const filename = `${crypto.randomUUID().split('-')[0]}${extname(file?.name)}`;

		proj.pages.push({
			title: title,
			src: filename,
			background: color
		});

		const buf = Buffer.from(await file.arrayBuffer());

		await s3.putObject(AWS_S3_DEFAULT_BUCKET, `${params.id}/project.json`, JSON.stringify(proj));
		await s3.putObject(AWS_S3_DEFAULT_BUCKET, `${params.id}/${filename}`, buf);
	}
} as Actions;
