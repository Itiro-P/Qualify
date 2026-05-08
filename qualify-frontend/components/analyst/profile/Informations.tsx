import { Briefcase, MapPin, Star } from "lucide-react";

export function Informations() {
  return (
    <div className="w-65/100 flex flex-col content-between">
      <h1 className="m-6 text-3xl">Testerson testersiano</h1>
      <p className="m-6 text-xl text-blue-800">Testador de teste</p>
      <div className="flex m-6 justify-start">
        <div className="flex">
          <MapPin className="text-accent mr-2" />
          <p className="mr-12 text-zinc-400">Testersian</p>
        </div>
        <div className="flex">
          <Briefcase className="text-accent mr-2" />
          <p className="mr-12 text-zinc-400">8+ semanas de experiência</p>
        </div>
        <div className="flex">
          <Star className="text-yellow-500 fill-current mr-2" />
          <p className="mr-12 text-zinc-400">
            4.9
            <span className="text-zinc-600 ml-1">20 reviews</span>
          </p>
        </div>
      </div>
    </div>
  );
}
