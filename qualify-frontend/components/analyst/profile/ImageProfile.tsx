import Image, { StaticImageData } from 'next/image';
import Testerson from "@/public/Testerson.png";

export function ImageProfile({
  user,
  imageURL,
}: {
  user: string;
  imageURL: string | null | undefined; // Deixando mais seguro caso a API mude
}) {
  
  // Forçamos o tipo do src a aceitar tanto a string da URL quanto o objeto estático do Next
  const imageSrc: string | StaticImageData = 
    imageURL && imageURL.trim() !== "" 
      ? `http://localhost:3001/uploads/${imageURL}` 
      : Testerson;

  return (
    <div className="w-[15%]">
      <Image
        src={imageSrc}
        alt={`Foto de ${user}`}
        className="inline-block object-cover" // Ajustado o size para o padrão Tailwind
        width={300}
        height={300}
        // Esse atributo avisa o Next para não tentar adivinhar as dimensões se for string externa
        unoptimized={typeof imageSrc === 'string'} 
      />
    </div>
  );
}