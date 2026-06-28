import Hero from "../components/Hero";
import Video from "../components/Video";
import About from "../components/About";
import Menu from "../components/Menu";
import Pricing from "../components/Pricing";

export default function HomePage() {
  return (
    <main>
      <Hero />
      <Video />
      <About />
      <Menu />
      {/* <Pricing /> */}
    </main>
  );
}
