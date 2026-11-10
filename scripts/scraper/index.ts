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

    // テーブルの描画を待ってから、テーブル上で右クリックしてメニューを開く
    const tableHeader = page.getByText("Page path", { exact: true }).first();
    await tableHeader.waitFor({ state: "visible", timeout: 60000 });
    await tableHeader.click({ button: "right" });

    try {
      await page.locator("text=Export Graph").first().click({ timeout: 6000 });
    } catch (err) {
      await page.locator("text=グラフをエクスポート").first().click({ timeout: 6000 });
    }

    try {
      await page.locator("text=Export Data").first().click({ timeout: 6000 });
    } catch (err) {
      await page.locator("text=データのエクスポート").first().click({ timeout: 6000 });
    }

    // 「エクスポート」を押してからダウンロード開始まで 20 秒以上かかることがある
    const [download] = await Promise.all([
      page.waitForEvent("download", { timeout: 60000 }),
      page.getByRole("button", { name: /^(Export|エクスポート)$/ }).click(),
    ]);

    await download.saveAs(savePath);
  } finally {
    await browser.close();
  }
};

scrape();
