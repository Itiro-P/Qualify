"use client";

import { useState } from "react";
import { Review } from "./Review";
import { Biografia } from "./Biografia";
import { Tecnologias } from "./Tecnologias";
import { Certificacoes } from "./Certificacoes";

function SistemaAbas() {
    const [abaAtiva, setAbaAtiva] = useState<"biografia" | "reviews" | "tecnologias" | "certificacoes">("biografia");

    return (
        <div>
            <div className="flex">
                <button className="p-2 m-2 bg-blue-950" onClick={() => setAbaAtiva("biografia")}>
                    Biografia
                </button>
                <button className="p-2 m-2 bg-blue-950" onClick={() => setAbaAtiva("reviews")}>
                    Reviews
                </button>
                <button className="p-2 m-2 bg-blue-950" onClick={() => setAbaAtiva("tecnologias")}>
                    Tecnologias
                </button>
                <button className="p-2 m-2 bg-blue-950" onClick={() => setAbaAtiva("certificacoes")}>
                    Certificações
                </button>
            </div>

            <div className="mt-4">
                {abaAtiva === "biografia" && (
                    <Biografia />
                )}

                {abaAtiva === "reviews" && (
                    <Review />
                )}

                {abaAtiva === "tecnologias" && (
                    <Tecnologias />
                )}

                {abaAtiva === "certificacoes" && (
                    <Certificacoes />
                )}
            </div>
        </div>
    );
}

export function Sobre(){
    return(
        <section id="sobre" className="flex flex-col px-3 w-8/10">
            <SistemaAbas />
        </section>
    );
}