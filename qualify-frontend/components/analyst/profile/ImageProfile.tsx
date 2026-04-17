import Image from 'next/image'
import photo from '@/public/Testerson.png';

var user: string = "Testerson";

export function ImageProfile(){
    return(
        <div className="w-15/100">
            <Image 
                src={photo} 
                alt={"Foto de "+{user}+" usuário"} 
                className="size-60 inline-block" 
            />
        </div>
    );
}