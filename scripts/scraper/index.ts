import path from "path";
import { fileURLToPath } from "url";

import { chromium } from "playwright";

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
    console.log("Page loaded");

    const mainArea = page.locator("body");
    const box = await mainArea.boundingBox();

    if (!box) {
      throw new Error(
        "Bounding box not found. Ensure the page is loaded correctly."
      );
    }

    const centerX = box.x + box.width / 2;
    const centerY = box.y + box.height / 2;
    await page.mouse.click(centerX, centerY, { button: "right" });
    console.log("Right-clicked");

    try {
      await page.locator("text=Export Data").first().click();
      console.log("Clicked context menu export");
    } catch (err) {
      console.log("Trying Japanese text...");
      await page.locator("text=データのエクスポート").first().click();
      console.log("Clicked context menu export (Japanese)");
    }

    await page.locator("text=エクスポート").nth(1).click();
    console.log("Clicked submenu export");

    const [download] = await Promise.all([
      page.waitForEvent("download", { timeout: 15000 }),
      page
        .locator('button:has-text("エクスポート"), button:has-text("Export")')
        .click(),
    ]);

    await download.saveAs(savePath);
    console.log("Download completed");
  } catch (err: unknown) {
    console.error("Error:", (err as Error).message);
  } finally {
    await browser.close();
  }
};

scrape();
