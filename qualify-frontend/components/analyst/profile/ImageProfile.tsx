import Image from "next/image";
import Testerson from "@/assets/images/Testerson.png";

export function ImageProfile({
  user,
  imageURL,
}: {
  user: string;
  imageURL: string;
}) {
  return (
    <div className="w-15/100">
      <Image
        src={imageURL == "" ? Testerson : imageURL}
        alt={"Foto de " + user + " usuário"}
        className="size-60 inline-block"
      />
    </div>
  );
}
