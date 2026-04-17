import Image from 'next/image'
import photo from "@/public/Testerson.png";
import type { ReactNode } from "react";

var user: string = "Testerson";
var imageProfile: ReactNode = <Image src="{photo}" alt={"Foto de "+{user}+" usuário"} className="size-60 inline-block" />;

export function ImageProfile(){
    return(
        <div className="w-15/100">
            {imageProfile}
        </div>
    );
}