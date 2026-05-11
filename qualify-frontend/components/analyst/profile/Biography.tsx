import { ReceiptText } from "lucide-react";

function BiographyCard({ biography }: { biography: string }) {
  return <p className="text-zinc-400 px-2 py-1 m-1 ">{biography}</p>;
}

export function Biography({ biography }: { biography: string }) {
  return (
    <div id="biography" className="flex flex-col p-4 mb-6">
      <div className="flex">
        <ReceiptText className="text-blue-500 mr-2" />
        <h2 className="text-xl">Biografia</h2>
      </div>
      <BiographyCard biography={biography} />
    </div>
  );
}
