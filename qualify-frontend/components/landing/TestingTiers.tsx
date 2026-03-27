import { SearchCheck, Settings2, ShieldCheck, Check } from "lucide-react";
import type { ReactNode } from "react";

interface TierCardProps {
  icon: ReactNode;
  title: string;
  description: string;
  features: string[];
}

function TierCard({ icon, title, description, features }: TierCardProps) {
  return (
    <div className="group p-8 rounded-xl border border-white/10 bg-white/5 hover:border-accent/50 transition-all">
      <div className="text-accent mb-6 group-hover:scale-110 transition-transform inline-block">
        {icon}
      </div>
      <h3 className="text-xl font-bold mb-3">{title}</h3>
      <p className="text-sm text-neutral-slate leading-relaxed mb-6">
        {description}
      </p>
      <ul className="text-xs space-y-2 text-slate-300">
        {features.map((feature) => (
          <li key={feature} className="flex items-center gap-2">
            <Check className="w-3.5 h-3.5 text-accent" />
            {feature}
          </li>
        ))}
      </ul>
    </div>
  );
}

const tiers: TierCardProps[] = [
  {
    icon: <SearchCheck className="w-10 h-10" />,
    title: "Teste de experiência",
    description:
      "Descubra casos extremos com intuição humana e cartas estruturadas. Perfeito para lançamentos de novas funcionalidades.",
    features: ["Auditoria de Fluxo UX/UI", "Descoberta de Casos Extremos"],
  },
  {
    icon: <Settings2 className="w-10 h-10" />,
    title: "Engenharia de Automação",
    description:
      "Construa pipelines robustos de CI/CD com Selenium, Playwright e Appium. Reduza o trabalho manual.",
    features: ["Configuração de Framework Personalizado", "Scripting de Regressão"],
  },
  {
    icon: <ShieldCheck className="w-10 h-10" />,
    title: "Auditorias de Segurança",
    description:
      "Identifique vulnerabilidades com testes de penetração abrangentes e verificações de conformidade SOC2.",
    features: ["Varredura de Vulnerabilidades", "Preparação para Conformidade"],
  },
];

export function TestingTiers() {
  return (
    <section id="tiers" className="px-6 md:px-20 py-24 bg-white/[0.02]">
      <div className="max-w-7xl mx-auto">
        <div className="mb-16">
          <h2 className="text-3xl font-bold mb-4">
            Tiers de Testes Especializados
          </h2>
          <div className="h-1 w-20 bg-accent mb-6" />
          <p className="text-neutral-slate max-w-2xl">
            Soluções personalizadas para cada etapa do seu ciclo de desenvolvimento de software, desde a auditoria inicial até a entrega contínua.
          </p>
        </div>
        <div className="grid md:grid-cols-3 gap-8">
          {tiers.map((tier) => (
            <TierCard key={tier.title} {...tier} />
          ))}
        </div>
      </div>
    </section>
  );
}
