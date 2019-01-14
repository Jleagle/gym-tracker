const puppeteer = require('puppeteer');

(async () => {
    const browser = await puppeteer.launch(
        {
            headless: true,
            args: [
                '--no-sandbox',
                // '--disable-setuid-sandbox',
                // '--enable-logging',
                // '--v=1'
            ]
        }
    );
    const page = await browser.newPage();
    await page.goto('https://www.puregym.com/login/');

    await page.type('#email', 'x');
    await page.type('#pin', 'y');
    await page.click('#login-submit');
    await page.waitForNavigation({
        timeout: 3000,
    });

    browser.close();

    console.log('Done');

})();
