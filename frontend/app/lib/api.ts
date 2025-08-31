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

export async function parseFiles(fileList: FileList) {}
