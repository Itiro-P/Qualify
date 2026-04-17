import { ImageProfile } from '@/components/analyst/profile/ImageProfile';
import { StatisticProfile } from "@/components/analyst/profile/StatisticProfile";
import { About } from "@/components/analyst/profile/About";
import { Informations } from "@/components/analyst/profile/Informations";
import { ContactButtons } from "@/components/analyst/profile/ContactButtons";

export function Profile(){
    return(
        <section id="profile" className="px-6 md:px-20 py-14">
            <div className="flex justify-between ml-3">
                <ImageProfile />
                <Informations />
                <ContactButtons />
            </div>
            <div className="flex mt-10">
                <About />
                <StatisticProfile />
            </div>
        </section>
    );
}