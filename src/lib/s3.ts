import {
	AWS_ENDPOINT_URL,
	AWS_ACCESS_KEY_ID,
	AWS_DEFAULT_REGION,
	AWS_SECRET_ACCESS_KEY
} from '$env/static/private';

import * as Minio from 'minio';

const client = new Minio.Client({
	endPoint: AWS_ENDPOINT_URL.split(':')[0],
	port: Number(AWS_ENDPOINT_URL.split(':')[1]),
	useSSL: false,
	region: AWS_DEFAULT_REGION,
	accessKey: AWS_ACCESS_KEY_ID,
	secretKey: AWS_SECRET_ACCESS_KEY
});

export default client;
