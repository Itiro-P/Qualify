import Image from 'next/image'
import foto from "@/public/Testerson.jpg";

export function ImagemPerfil(){
    return(
        <div className="w-15/100">
            <Image 
                src={foto} 
                alt="Foto Testerson"
                className="size-60 inline-block"
            />
        </div>
    );
}