"use client";
import UploadButton from "@/app/_components/upload-button";
import RenewalCalendar from "@/app/_components/renewal-calendar";
import { useState } from "react";

export default function Home() {
  const renewalEvents = [
    {
      vendor: "Vendor A",
      effectiveDate: "2023-01-01",
      renewalDate: "2025-08-29",
      autoRenew: true,
      amount: 1000,
    },
    {
      vendor: "Vendor B",
      effectiveDate: "2023-02-01",
      renewalDate: "2025-09-01",
      autoRenew: false,
      amount: 2000,
    },
  ];

  async function handleButtonClick() {
    console.log("Button clicked");
    const res = await fetch("https://openrouter.ai/api/v1/chat/completions", {
      method: "POST",
      headers: {
        Authorization:
          "Bearer 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",
        "HTTP-Referer": "<YOUR_SITE_URL>", // Optional. Site URL for rankings on openrouter.ai.
        "X-Title": "<YOUR_SITE_NAME>", // Optional. Site title for rankings on openrouter.ai.
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        model: "openai/gpt-4o",
        messages: [
          {
            role: "user",
            content: "What is the meaning of life?",
          },
        ],
      }),
    });

    // Parse the JSON response
    const data = await res.json();

    // The API returns an array of choices
    console.log("Full response:", data);
  }

  return (
    <>
      <h1 className="mt-8 mb-4 text-center text-3xl font-bold">
        Renewal Calendar
      </h1>
      <div className="flex flex-col items-center justify-center min-h-[10vh]">
        <UploadButton />
      </div>

      <div>
        <RenewalCalendar renewalEvents={renewalEvents} />
      </div>

      <button onClick={handleButtonClick}>Test OpenRouter API</button>
    </>
  );
}
