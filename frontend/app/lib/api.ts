import { VendorContract } from "./definitions";

/**
 * Uploads a list of files to S3.
 * @param fileList The list of files to upload.
 */

const baseURL =
  process.env.NEXT_PUBLIC_GATEWAY_API_URL || "http://localhost:8080";
export async function uploadToS3(fileList: FileList) {
  const formData = new FormData();
  Array.from(fileList).forEach((file) => {
    formData.append("file", file);
  });

  try {
    const response = await fetch(`${baseURL}/upload`, {
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

/**
 * Parses the uploaded PDF files to extract seller, effective date, renewal date,
 * and whether contract autorenews. Handled by backend API and uploads to database.
 * @param fileList The list of files to parse.
 */
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
    const response = await fetch(`${baseURL}/parse`, {
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

/**
 * Converts a byte array to a base64 string.
 * @param bytes The byte array to convert.
 * @returns The base64-encoded string.
 */
function bytesToBase64(bytes: Uint8Array) {
  let binary = "";
  for (let i = 0; i < bytes.length; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary);
}

/**
 * Fetches calendar data from the backend API to display.
 * @returns An array of VendorContract objects.
 */
export async function fetchCalendarData() {
  try {
    const response = await fetch(`${baseURL}/calendarData`, {
      method: "GET",
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(
        `Failed to fetch calendar data ${response.status} ${errorText}`,
      );
    }

    const data: VendorContract[] = await response.json();
    console.log("Fetched calendar data:", data);
    return data;
  } catch (error: any) {
    console.log("Error fetching calendar data:", (error as Error).message);
    return [];
  }
}
