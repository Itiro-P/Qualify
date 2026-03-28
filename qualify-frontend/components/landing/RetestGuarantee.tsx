import { ShieldCheck, History, Banknote, Handshake } from "lucide-react";

export function RetestGuarantee() {
  return (
    <section id="guarantee" className="px-6 md:px-20 py-24 bg-primary/10">
      <div className="max-w-5xl mx-auto glass-panel p-12 rounded-2xl border-white/20 text-center relative overflow-hidden">
        <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-accent to-transparent" />
        <ShieldCheck className="w-16 h-16 text-accent mx-auto mb-6" />
        <h2 className="text-4xl font-bold mb-6">Garantia de Reteste</h2>
        <p className="text-neutral-slate text-xl mb-10 max-w-2xl mx-auto leading-relaxed">
          Se um bug que perdemos chegar à produção, ou se uma correção não for
          verificada de acordo com sua satisfação, fornecemos reteste imediato
          em 48 horas{" "}
          <span className="text-white font-bold">sem custo adicional.</span>
        </p>
        <div className="flex flex-wrap justify-center gap-8">
          <div className="flex items-center gap-3">
            <History className="w-6 h-6 text-accent" />
            <span className="text-sm font-medium">48h de reteste</span>
          </div>
          <div className="flex items-center gap-3">
            <Banknote className="w-6 h-6 text-accent" />
            <span className="text-sm font-medium">Sem taxas ocultas</span>
          </div>
          <div className="flex items-center gap-3">
            <Handshake className="w-6 h-6 text-accent" />
            <span className="text-sm font-medium">SLA Protegido</span>
          </div>
        </div>
      </div>
    </section>
  );
}
