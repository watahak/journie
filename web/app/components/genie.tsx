export default function Genie({ size }: { size: "sm" | "lg" }) {
  const width = size === "sm" ? "240" : "360";

  return (
    <video width={width} height="240" preload="none" loop={true} autoPlay muted playsInline>
      <source src="https://mt5gkfie7qymwe4f.public.blob.vercel-storage.com/genie-lwmVPIo3crhQ0ldlrSkHgoS9E79IAd.mov" type="video/mp4" />
      Your browser does not support the video tag.
    </video>
  );
}
