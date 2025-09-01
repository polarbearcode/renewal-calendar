"use server";
export async function uploadToS3(fileList: FileList) {
  const formData = new FormData();
  Array.from(fileList).forEach((file) => {
    formData.append("file", file);
  });

  try {
    const response = await fetch("http://localhost:8080/upload", {
      method: "POST",
      body: formData,
      headers: {
        "X-Upload-Timestamp": new Date().toISOString(),
      },
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

  const fileMap = new Map<File, Uint8Array>();

  for (const file of fileList) {
    const arrayBuffer = await file.arrayBuffer();
    fileMap.set(file, new Uint8Array(arrayBuffer));
  }

  const fileMapObj = Object.fromEntries(fileMap);
  const fileMapJSON = JSON.stringify(fileMapObj);

  try {
    const response = await fetch("http://localhost:8080/parse", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: fileMapJSON,
    });
    if (!response.ok) {
      throw new Error("Failed to parse files");
    }

    const data = await response.json();
    console.log("Parsed result:", data.result);
  } catch (error) {
    console.log("Error with parsing:", error);
  }
}
