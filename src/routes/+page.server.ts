import { fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

import { db, s3, type Project } from '$lib';
import { AWS_S3_DEFAULT_BUCKET } from '$env/static/private';

export const load = (async ({}) => {
	const res = await db.all<Project[]>('SELECT ID, Name FROM projects');
	return { projects: res };
}) as PageServerLoad;

export const actions: Actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		const name = data.get('project-name');

		if (!name) return fail(400, { name, missing: true });

		const uuid = crypto.randomUUID().split('-')[0];

		const res = await db.run('INSERT OR IGNORE INTO projects (ID, Name) VALUES (:id, :name)', {
			':id': uuid,
			':name': name
		});

		const project: Project = {
			id: uuid,
			title: name.toString(),
			pages: []
		};

		await s3.putObject(AWS_S3_DEFAULT_BUCKET, `${uuid}/project.json`, JSON.stringify(project));

		if (res.changes == undefined) {
			return fail(500, { reason: 'Failed to insert project into database' });
		}

		return { success: true };
	}
};
