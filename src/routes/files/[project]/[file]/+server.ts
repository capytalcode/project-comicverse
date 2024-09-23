import { type RequestHandler } from '@sveltejs/kit';
import stream from 'node:stream/promises';

import { db, s3, type Project } from '$lib';
import { AWS_S3_DEFAULT_BUCKET } from '$env/static/private';
import { extname } from 'node:path';

export const GET = (async ({ params }) => {
	const file = await s3.getObject(AWS_S3_DEFAULT_BUCKET, `${params.project}/${params.file}`);
	file.on('error', (err: any) => {
		console.log(err);
	});

	let chunks: Buffer[] = [];
	let buf;

	file.on('data', (chunk) => {
		chunks.push(Buffer.from(chunk));
	});
	file.on('end', () => {
		buf = Buffer.concat(chunks);
	});
	await stream.finished(file)

	let res = new Response(buf);
	res.headers.set(
		'Content-Type',
		(() => {
			switch (extname(params.file!)) {
				case '.png':
					return 'image/png';
				case '.json':
					return 'application/json';
			}
			return 'text/plain';
		})()
	);
	res.headers.set('Cache-Control', 'max-age=604800')

	return res;
}) as RequestHandler;
