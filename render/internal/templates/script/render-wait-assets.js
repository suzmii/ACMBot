() => {
	const waitForNextFrames = (count = 2) =>
		new Promise(resolve => {
			const step = () => {
				if (count <= 0) {
					resolve();
					return;
				}
				count -= 1;
				requestAnimationFrame(step);
			};
			requestAnimationFrame(step);
		});

	const waitForImage = (img, timeoutMs = 5000) =>
		new Promise(resolve => {
			let settled = false;
			const finish = () => {
				if (settled) {
					return;
				}
				settled = true;
				resolve();
			};

			const timer = setTimeout(finish, timeoutMs);
			const cleanup = () => {
				clearTimeout(timer);
				img.removeEventListener('load', onDone);
				img.removeEventListener('error', onDone);
			};
			const onDone = () => {
				cleanup();
				if (typeof img.decode === 'function') {
					img.decode().catch(() => {}).finally(finish);
					return;
				}
				finish();
			};

			if (img.complete) {
				onDone();
				return;
			}

			img.addEventListener('load', onDone, { once: true });
			img.addEventListener('error', onDone, { once: true });
		});

	const backgroundURLs = new Set();
	for (const el of document.querySelectorAll('*')) {
		const backgroundImage = getComputedStyle(el).backgroundImage;
		if (!backgroundImage || backgroundImage === 'none') {
			continue;
		}
		for (const match of backgroundImage.matchAll(/url\((['"]?)(.*?)\1\)/g)) {
			if (match[2]) {
				backgroundURLs.add(match[2]);
			}
		}
	}

	const waitForBackgroundImage = (src, timeoutMs = 5000) =>
		new Promise(resolve => {
			const img = new Image();
			const finish = () => {
				clearTimeout(timer);
				resolve();
			};
			const timer = setTimeout(finish, timeoutMs);
			img.onload = () => {
				if (typeof img.decode === 'function') {
					img.decode().catch(() => {}).finally(finish);
					return;
				}
				finish();
			};
			img.onerror = finish;
			img.src = src;
			if (img.complete) {
				img.onload();
			}
		});

	return Promise.resolve()
		.then(() => document.fonts.ready)
		.then(() =>
			Promise.all([
				...Array.from(document.images, img => waitForImage(img)),
				...Array.from(backgroundURLs, src => waitForBackgroundImage(src)),
			]),
		)
		.then(() => waitForNextFrames(2));
};
