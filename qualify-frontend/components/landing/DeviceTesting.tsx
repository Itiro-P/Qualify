import Image from "next/image";
import { Smartphone, Watch, Tv, Router } from "lucide-react";
import type { ReactNode } from "react";

interface DeviceCategoryProps {
  icon: ReactNode;
  title: string;
  subtitle: string;
}

function DeviceCategory({ icon, title, subtitle }: DeviceCategoryProps) {
  return (
    <div className="flex items-start gap-3">
      <span className="text-accent">{icon}</span>
      <div>
        <h4 className="font-bold text-white">{title}</h4>
        <p className="text-xs text-neutral-slate">{subtitle}</p>
      </div>
    </div>
  );
}

const deviceImages = [
  {
    src: "https://lh3.googleusercontent.com/aida-public/AB6AXuAkRQu5tf6spLnToJpdM4XeYQwg-aD7eYOrbyXb0NmkZxtZt2LxlOUYbXOljBiiMzAWy4YnorGTAkMUM0Sjp7aEg-snDt4cmJSrRwGHKo-FnuMNqocIk_8yV8zckAu3BQW4qCnKMMZ2VohjMc5AZ5AmzeBSR8aGmBQUEBfUePqFZcODISBlPZLGMMftnO4VFboYOydyPtCJf1iZAWqne_bkFdu6FvmF9ch_2-yf9scew79a08ZeBTG3QmHTFYRIPH9NfcFLQXBIz3tH",
    alt: "Selection of modern smartphones on a desk",
    height: "h-48",
  },
  {
    src: "https://lh3.googleusercontent.com/aida-public/AB6AXuBRBY2PKDseek7VfqjooBLtPLd_c9TPxHorMfqxJODqWNHMZN3lSf_S0YKTUokNikI9NJaPkysxbIuQGR2Fh7KdfC3v8byVoCTF1SXkfjql6dKBwaZt-TA0LGgpCDNFj-zui7b3qAP-O-BOeuIiZrpSFPq3LgDBRsVn1Ny3VCwUqPtPesxiqtphb9ynDyYyFTtaR3-nM1M3bcDAHg2W5KkWalZraxFGvXuoswPxYzJru1AupeRGJ_H4xymU5B818r7lXFQyWHQ9GhAx",
    alt: "High-end laptop displaying software interface",
    height: "h-64",
  },
  {
    src: "https://lh3.googleusercontent.com/aida-public/AB6AXuDKbg90gan6WStfhlsJyhWI9ncYEdJrxAWRmMPdcysvKEBGfh6-KEj-cJM7_HhfyyFWHEDYH11r6GoC6JCzHW3vd3F1gUpbRRAPo2kKrbQbTxg_kb8p5gmhuazVhVn6fjwC4B6eZbgLTkWiX_C5Z_pNibaUJTf19BVwit6dmsxU2rEakS4cY-Uytd7qawhQI1FxyweSLQ1rp4f8pPtL3y4OixsXVmVaCpsJcxwoCgmIYYBxzu1Ji-owd52BAFPTJcdxutaGl0jRgXf9",
    alt: "Professional gaming setup with multiple screens",
    height: "h-64",
  },
  {
    src: "https://lh3.googleusercontent.com/aida-public/AB6AXuAYyW0Zo5e_cZ-UH0CUs4RaIePjbNl-KYGVpiODpBa5PbbaPsjgL3Q1GsHpK0S26mtpRntUF_gu9hlE1dqX3G2DhbyaLbnjHG1tj3jIbPmqEq8RvAeij5h3jtf8wbYW2MtYogeeooX82BWIVqBiNHUGvjsbCI1OhA4z3Fcaqq2CmJ-arkkRvmX6fQo_9s9age65bQHf8sr-0amyalhugy-iUhQnaTIZOXx0QZnfXCn9EBI32pLzQxy_qVGKObwywhUeoQLTLbThdXMO",
    alt: "Smart home devices and IoT gadgets",
    height: "h-48",
  },
];

const categories: DeviceCategoryProps[] = [
  {
    icon: <Smartphone className="w-6 h-6" />,
    title: "Laboratórios Móveis",
    subtitle: "Variantes iOS & Android",
  },
  {
    icon: <Watch className="w-6 h-6" />,
    title: "Relógios Inteligentes",
    subtitle: "WatchOS & Tizen",
  },
  {
    icon: <Tv className="w-6 h-6" />,
    title: "Smart TV",
    subtitle: "Roku, FireTV, WebOS",
  },
  {
    icon: <Router className="w-6 h-6" />,
    title: "IoT / Embedded",
    subtitle: "Testes de firmware personalizado",
  },
];

export function DeviceTesting() {
  return (
    <section id="devices" className="px-6 md:px-20 py-24">
      <div className="max-w-7xl mx-auto flex flex-col lg:flex-row gap-20 items-center">
        <div className="lg:w-1/2 grid grid-cols-2 gap-4">
          <div className="space-y-4">
            <div
              className={`${deviceImages[0].height} rounded-xl overflow-hidden border border-white/10 relative`}
            >
              <Image
                className="object-cover"
                alt={deviceImages[0].alt}
                src={deviceImages[0].src}
                fill
                sizes="(max-width: 768px) 50vw, 25vw"
              />
            </div>
            <div
              className={`${deviceImages[1].height} rounded-xl overflow-hidden border border-white/10 relative`}
            >
              <Image
                className="object-cover"
                alt={deviceImages[1].alt}
                src={deviceImages[1].src}
                fill
                sizes="(max-width: 768px) 50vw, 25vw"
              />
            </div>
          </div>
          <div className="space-y-4 pt-8">
            <div
              className={`${deviceImages[2].height} rounded-xl overflow-hidden border border-white/10 relative`}
            >
              <Image
                className="object-cover"
                alt={deviceImages[2].alt}
                src={deviceImages[2].src}
                fill
                sizes="(max-width: 768px) 50vw, 25vw"
              />
            </div>
            <div
              className={`${deviceImages[3].height} rounded-xl overflow-hidden border border-white/10 relative`}
            >
              <Image
                className="object-cover"
                alt={deviceImages[3].alt}
                src={deviceImages[3].src}
                fill
                sizes="(max-width: 768px) 50vw, 25vw"
              />
            </div>
          </div>
        </div>
        <div className="lg:w-1/2">
          <h2 className="text-3xl font-bold mb-6">
            Testes em Dispositivos Reais
          </h2>
          <p className="text-neutral-slate text-lg mb-8 leading-relaxed">
            A simulação nunca é suficiente. Nossos especialistas testam em uma
            vasta gama de dispositivos reais para garantir que seu software
            funcione perfeitamente em qualquer ambiente.
          </p>
          <div className="grid grid-cols-2 gap-6 mb-8">
            {categories.map((category) => (
              <DeviceCategory key={category.title} {...category} />
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
