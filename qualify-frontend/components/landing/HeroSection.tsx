import Image from "next/image";
import { ArrowRight } from "lucide-react";

export function HeroSection() {
  return (
    <section className="relative px-6 md:px-20 pt-20 pb-32 overflow-hidden">
      <div className="absolute top-0 right-0 -z-10 w-1/2 h-full bg-gradient-to-l from-primary/10 to-transparent" />
      <div className="max-w-7xl mx-auto grid lg:grid-cols-2 gap-16 items-center">
        <div className="flex flex-col gap-8">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-accent/10 border border-accent/20 text-accent text-xs font-bold uppercase tracking-widest w-fit">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-accent opacity-75" />
              <span className="relative inline-flex rounded-full h-2 w-2 bg-accent" />
            </span>
            Apenas especialistas verificados
          </div>
          <h1 className="text-5xl md:text-7xl font-bold leading-[1.1] tracking-tight">
            Sem erros.
            <br />
            <span className="text-accent">Entrega rápida.</span>
          </h1>
          <p className="text-lg text-neutral-slate max-w-lg leading-relaxed">
            Conecte-se com os profissionais de QA mais experientes do mundo para
            testes exploratórios, automação e segurança. Implemente com 100%
            de confiança.
          </p>
          <div className="flex flex-wrap gap-4">
            <button className="bg-primary hover:scale-105 transition-transform px-8 py-4 font-bold rounded-lg flex items-center gap-2 cursor-pointer">
              Contratar especialistas
              <ArrowRight className="w-5 h-5" />
            </button>
            <button className="bg-white/5 hover:bg-white/10 border border-white/10 px-8 py-4 font-bold rounded-lg cursor-pointer">
              Explorar o marketplace
            </button>
          </div>
        </div>
        <div className="relative">
          <div className="aspect-square rounded-2xl overflow-hidden glass-panel p-2">
            <Image
              className="w-full h-full object-cover rounded-xl opacity-80"
              alt="Código técnico em um monitor escuro com acentos neon"
              src="https://lh3.googleusercontent.com/aida-public/AB6AXuDAgwZGP6X2i5Aq4GD_8rlJldPPnFEn0Byq6fF-W-gCObRdRoXzhE1RKXr8LQiMYweHLuczhdzYQLMWa4AFJC4h4ijXOJv01FdtqdNsPM72otwB4dCm1iwaewPfv3ZcQAPOclFnz3pKVhBuaoP0Uv7YAJCr32jj6hBegzWPB2BZ8HNqY6DqJV4FJn-w7gJOqYei7e5s4j3r-BWviwNdmDH_31rZJWHxiwLSv52olgIg_-wMW1mlTH-TJKo9zpESph8_sVnUryNmZdPC"
              width={600}
              height={600}
            />
          </div>
          {/* Floating metrics */}
          <div className="absolute -bottom-6 -left-6 glass-panel p-6 rounded-xl border-accent/30 mint-glow">
            <p className="text-accent font-bold text-3xl">2.5k+</p> {/* TODO: fetch real data */}
            <p className="text-xs uppercase tracking-wider text-neutral-slate">
              Especialistas ativos
            </p>
          </div>
          <div className="absolute -top-6 -right-6 glass-panel p-6 rounded-xl border-primary/30">
            <p className="text-white font-bold text-3xl">99.9%</p>
            <p className="text-xs uppercase tracking-wider text-neutral-slate">
              Satisfação do cliente
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
