document.querySelector('#accent-color-hue').addEventListener('change', (e) => {
	// eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
	const hue = /** @type {number} */ (e.target?.value);

	document.querySelector('body')
		.setAttribute('style', `--user-theme-accent-hue:${String(hue)};`);
});
