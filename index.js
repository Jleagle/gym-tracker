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

    await page.type('#email', process.env.PUREGYM_EMAIL);
    await page.type('#pin', process.env.PUREGYM_PIN);
    await page.click('#login-submit');

    try {

        await page.waitForNavigation({
            timeout: 3000,
        });

    } catch (e) {
        console.log('Login failed');
        console.log(e);
    }

    const people = await page.$eval('.heading.heading--level3.secondary-color.margin-none', el => el.innerText);

    console.log(people);

    browser.close();
    console.log('Done');

})();
