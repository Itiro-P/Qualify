import { ImagemPerfil } from './ImagemPerfil';
import { AoLado } from "./AoLado";
import { Sobre } from "./Sobre";
import { Info } from "./Info";
import { BotoesContato } from "./BotoesContato";

export function Perfil(){
    return(
        <section id="perfil" className="px-6 md:px-20 py-14">
            <div className="flex justify-between ml-3">
                <ImagemPerfil />
                <Info />
                <BotoesContato />
            </div>
            <div className="flex mt-10">
                <Sobre />
                <AoLado />
            </div>
        </section>
    );
}