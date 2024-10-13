import {
	defineConfig,
	presetIcons,
	presetTypography,
	presetUno,
	presetWebFonts,
	transformerDirectives,
	transformerVariantGroup,
} from 'unocss';

export default defineConfig({
	cli: {
		entry: {
			outFile: './assets/css/uno.css',
			patterns: [
				'./{pages,templates}/**/*.templ',
				'./assets/**/*.{js,css,html}',
				'!./assets/css/uno.css',
			],
		},
	},
	presets: [
		presetIcons(),
		presetTypography(),
		presetUno(),
		presetWebFonts({
			fonts: {
				cal: {
					name: 'Cal Sans',
					provider: 'none',
				},
				mono: {
					name: 'Fira Code',
					provider: 'none',
				},
				sans: {
					name: 'Inter',
					provider: 'none',
				},
			},
			provider: 'bunny',
		}),
	],
	transformers: [transformerDirectives(), transformerVariantGroup()],
});
