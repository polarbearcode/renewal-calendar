import { VendorContract } from "./definitions";

export async function uploadToS3(fileList: FileList) {
  console.log("xyz");
  const formData = new FormData();
  Array.from(fileList).forEach((file) => {
    formData.append("file", file);
  });

  try {
    const response = await fetch("http://localhost:8080/upload", {
      method: "POST",
      body: formData,
    });
    if (!response.ok) {
      throw new Error("Failed to upload files");
    }
  } catch (error) {
    console.error("Error uploading files:", error);
  }
}

export async function parseFiles(fileList: FileList) {
  // make a map of files -> bytes

  const fileMap = new Map<string, string>();

  for (const file of fileList) {
    const arrayBuffer = await file.arrayBuffer();
    const bytes = new Uint8Array(arrayBuffer);
    fileMap.set(file.name, bytesToBase64(bytes));
  }

  const fileMapObj = Object.fromEntries(fileMap);

  const payload = {
    fileBytesMap: fileMapObj,
  };

  const fileMapJSON = JSON.stringify(payload);

  try {
    const response = await fetch("http://localhost:8080/parse", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: fileMapJSON,
    });
    if (!response.ok) {
      console.log(response.status);
      throw new Error("Failed to parse files");
    }

    const data = await response.json();
    console.log("Parsed result:", data.result);
  } catch (error) {
    console.log("Error with parsing:", error);
  }
}

function bytesToBase64(bytes: Uint8Array) {
  let binary = "";
  for (let i = 0; i < bytes.length; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary);
}

export async function fetchCalendarData() {
  try {
    const response = await fetch("http://localhost:8080/calendarData", {
      method: "GET",
    });
    if (!response.ok) {
      throw new Error("Failed to fetch calendar data");
    }

    const data: VendorContract[] = await response.json();
    console.log("Fetched calendar data:", data);
    return data;
  } catch (error) {
    console.log("Error fetching calendar data:", error);
    return [];
  }
}
