import {
    ImageProfile,
    StatisticProfile,
    About,
    Informations,
    ContactButtons
} from '@/components/analyst/profile';

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