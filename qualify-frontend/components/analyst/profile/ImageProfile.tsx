import Image from "next/image";

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
        src={imageURL}
        alt={"Foto de " + user + " usuário"}
        className="size-60 inline-block"
      />
    </div>
  );
}
