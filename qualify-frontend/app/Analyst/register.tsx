import {
  RegisterAnalyst
} from "@/components/analyst/register";

import { Footer, Header } from "@/components";

export function Register() {
  return (
    <section id="register" className="px-6 md:px-20 py-14">
      <Header />
      <div className="flex justify-center">
        <div className="flex flex-col ml-3 w-1/2">
          <RegisterAnalyst />
        </div>
      </div>
      <Footer />
    </section>
  );
}
