import { ReceiptText } from 'lucide-react';

const biografia: string = "Meu nome é terterson testersiano, sou de testersian e tenho 8 semanas de experiência, em testes aleatórios e teste de botões, além do mais sou um guardião da qualidade."

export function BiografiaCard({biografia}: {biografia : string} ) {
    return (
        <p className="text-zinc-400 px-2 py-1 m-1 ">
            {biografia}
        </p>
    );
}

export function Biografia() {
    return (
        <div id="biografia" className="flex flex-col p-4 mb-6">
            <div className="flex">
                <ReceiptText className="text-blue-500 mr-2"/>
                <h2 className="text-xl">
                    Biografia
                </h2>
            </div>
            <BiografiaCard biografia={biografia} />
        </div>
    );
}