import { Briefcase, MapPin, Star } from "lucide-react";

function getTimeOnPlatform(createdDateTex: string) {
  const createdDate = new Date(createdDateTex);
  const now = new Date();

  const diffMs = now.getTime() - createdDate.getTime();

  const days = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  const years = Math.floor(days / 365);
  const months = Math.floor(days / 30);

  if (years > 0) {
    return `${years} ano${years > 1 ? "s" : ""}`;
  }

  if (months > 0) {
    return `${months} mês${months > 1 ? "es" : ""}`;
  }

  return `${days} dia${days > 1 ? "s" : ""}`;
}

export function Informations({
  name,
  city,
  rating,
  reviews,
  date,
}: {
  name: string;
  city: string;
  rating: number;
  reviews: number;
  date: string;
}) {
  return (
    <div className="w-[65%] flex flex-col content-between justify-between">
      <h1 className="m-6 text-3xl">{name}</h1>
      <div className="flex m-6 justify-start">
        <div className="flex">
          <MapPin className="text-accent mr-2" />
          <p className="mr-12 text-zinc-400">{city}</p>
        </div>
        <div className="flex">
          <Briefcase className="text-accent mr-2" />
          <p className="mr-12 text-zinc-400">
            {getTimeOnPlatform(date)} na plataforma
          </p>
        </div>
        <div className="flex">
          <Star className="text-yellow-500 fill-current mr-2" />
          <p className="mr-12 text-zinc-400">
            {rating}
            <span className="text-zinc-600 ml-1">{reviews} reviews</span>
          </p>
        </div>
      </div>
    </div>
  );
}
