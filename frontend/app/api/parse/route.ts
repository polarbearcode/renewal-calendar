import { NextRequest, NextResponse } from "next/server";
import fs from "fs";
import os from "os";
import path from "path";
import pdf from "pdf-parse";

export const runtime = "nodejs";

export async function POST(req: NextRequest) {
  let tempPath: string | null = null;

  try {
    console.log("POST handler invoked");

    // 1️⃣ Get the uploaded file
    const formData = await req.formData();
    const file = formData.get("file") as File;

    if (!file) {
      return NextResponse.json({ error: "No file uploaded" }, { status: 400 });
    }

    // 2️⃣ Convert to buffer and save to a temporary path
    const buffer = Buffer.from(await file.arrayBuffer());
    const safeName = path.basename(file.name); // removes any directories from client
    tempPath = path.join(os.tmpdir(), `${Date.now()}-${safeName}`);
    await fs.promises.writeFile(tempPath, buffer);

    // 3️⃣ Read PDF from temp path
    const pdfBuffer = await fs.promises.readFile(tempPath);
    const pdfData = await pdf(pdfBuffer);
    const text = pdfData.text;

    // 4️⃣ Send PDF text to OpenRouter
    const chatRes = await fetch(
      "https://openrouter.ai/api/v1/chat/completions",
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${process.env.OPENROUTER_API_KEY!}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          model: "gpt-4o-mini",
          messages: [
            {
              role: "user",
              content: `Extract the main key points from this PDF:\n\n${text}`,
            },
          ],
        }),
      },
    );

    if (!chatRes.ok) {
      const errText = await chatRes.text();
      console.error("Chat completion failed:", errText);
      return NextResponse.json({ error: errText }, { status: chatRes.status });
    }

    const chatData = await chatRes.json();
    const result = chatData.choices?.[0]?.message?.content ?? "No content";

    return NextResponse.json({ result });
  } catch (err) {
    console.error(err);
    return NextResponse.json(
      { error: (err as Error).message },
      { status: 500 },
    );
  } finally {
    // 5️⃣ Clean up temporary file
    if (tempPath) {
      try {
        await fs.promises.unlink(tempPath);
      } catch {
        // ignore cleanup errors
      }
    }
  }
}
