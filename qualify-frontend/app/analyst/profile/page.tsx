"use client";

import { Footer, Header } from "@/components";
import {
  ImageProfile,
  StatisticProfile,
  About,
  Informations,
  ContactButtons,
} from "@/components/analyst/profile";
import { useState } from "react";
import { Analyst } from "@/types/services";

export default function Profile() {
  const [analyst, setAnalyst] = useState<Analyst | null>({
    id: 0,
    name: "",
    email: "",
    phone: "",
    city: "",
    country_code: "",
    country_name: "",
    country_state: "",
    timezone: "",
    hourly_rate: 0,
    mean_rating: 0,
    total_reviews: 0,
    time_created: "",
  });

  //tem que implementar forma de pegar o usuario da seessão e ver se é analista ou não e colocar no 'analyst'

  return (
    <section id="profile" className="px-6 md:px-20 py-14">
      <Header />
      {analyst ? (
        <div>
          <div className="flex justify-between ml-3 mt-10">
            <ImageProfile />
            <Informations />
            <ContactButtons />
          </div>

          <div className="flex mt-10">
            <About analyst={analyst} />
            <StatisticProfile />
          </div>
        </div>
      ) : (
        <p>Usuário não logado ou não é analista</p>
      )}

      <Footer />
    </section>
  );
}
