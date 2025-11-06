import path from "path";
import { chromium } from "playwright";
import { fileURLToPath } from "url";

const url =
  "https://lookerstudio.google.com/u/0/reporting/fe615dc5-59b0-4db5-8504-ef9eacb663a9/page/p_ognozybded";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const savePath = path.resolve(__dirname, "../../data/master/page_view.csv");

const scrape = async () => {
  const browser = await chromium.launch({
    headless: true,
    args: ["--no-sandbox", "--disable-setuid-sandbox"],
  });
  const context = await browser.newContext({
    acceptDownloads: true,
    locale: "ja-JP",
    viewport: { width: 1280, height: 800 },
  });
  const page = await context.newPage();

  try {
    await page.goto(url);
    await page.waitForTimeout(10000);

    const mainArea = page.locator("body");
    const box = await mainArea.boundingBox();

    if (!box) {
      throw new Error("Bounding box not found");
    }

    const centerX = box.x + box.width / 2;
    const centerY = box.y + box.height / 2;
    await page.mouse.click(centerX, centerY, { button: "right" });

    try {
      await page.locator("text=Export Data").first().click();
    } catch (err) {
      await page.locator("text=データのエクスポート").first().click();
    }

    await page.locator("text=エクスポート").first().click();

    const [download] = await Promise.all([
      page.waitForEvent("download", { timeout: 15000 }),
      page.locator("text=エクスポート").first().click(),
    ]);

    await download.saveAs(savePath);
  } finally {
    await browser.close();
  }
};

scrape();