"use client";
import { useRef } from "react";

// Component for uploading PDF files button
export default function UploadButton() {
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files && files.length > 0) {
      console.log("Selected files:");
      Array.from(files).forEach((file) => console.log(file.name));
    }
  };

  return (
    <div>
      <button
        className="bg-blue-500 text-white p-2 rounded text-xl"
        onClick={handleButtonClick}
      >
        Upload PDF
      </button>
      <input
        type="file"
        multiple
        ref={fileInputRef}
        onChange={handleFileChange}
        accept="application/pdf"
        className="hidden"
      />
    </div>
  );
}
