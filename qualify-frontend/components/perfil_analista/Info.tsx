import Image from 'next/image'
import foto from "@/public/Testerson.jpg";
import { Briefcase, MapPin, Star } from 'lucide-react';

export function Info(){
    return(
        <section id="info" className="px-6 md:px-20 py-24">
            <div className="flex justify-between ">
                <div className="w-15/100">
                    <Image 
                        src={foto} 
                        alt="Foto Testerson"
                        className="size-60 inline-block"
                    />
                </div>
                <div className="w-65/100 flex flex-col content-between">
                    <h1 className="m-6 text-3xl">Testerson testersiano</h1>
                    <p className="m-6 text-xl text-blue-800">Testador de teste</p>
                    <div className="flex m-6 justify-start">
                        <div className="flex">
                            <MapPin className="text-accent mr-2"/>
                            <p className="mr-12 text-zinc-400">Testersian</p>
                        </div>
                        <div className="flex">
                            <Briefcase className="text-accent mr-2"/>
                            <p className="mr-12 text-zinc-400">8+ semanas de experiência</p>
                        </div>
                        <div className="flex">
                            <Star className="text-yellow-500 fill-current mr-2"/>
                            <p className="mr-12 text-zinc-400">4.9 <span className="text-zinc-600 ml-1">20 reviews</span></p>
                        </div>
                    </div>
                </div>
                <div className="w-20/100 flex items-end content-end">
                    <button className="px-5 py-3 m-5 bg-zinc-900 shadow-lg shadow-zinc-900/40 inset-ring-2 inset-ring-zinc-800">Contate</button>
                    <button className="px-5 py-3 m-5 bg-blue-700 shadow-lg shadow-blue-900/40 inset-ring-2 inset-ring-blue-600">Contrate agora</button>
                </div>
            </div>
            <div className="flex justify-between content-center">

            </div>
        </section>
    );
}