'use client'
import Image from 'next/image';
import Tilt from 'react-parallax-tilt';

export default function Cta({ size }: { size: "sm" | "lg" }) {
  const width = size === "sm" ? "240" : "360";

    return (
        <Tilt>
        <a href="tg://resolve?domain=Journie_Bot">
          <Image
            src="/qr.png"
            alt="Journie QR code"
            className=""
            width={width}
            height={width}
            priority
          />
        </a>
        </Tilt>
    )

}
