import { Play } from "lucide-react";
import { useState } from "react";

export default function Video() {
  const [isPlaying, setIsPlaying] = useState(false);

  return (
    <section className="px-4 py-20 sm:px-6 lg:px-8 bg-foreground">
      <div className="max-w-6xl mx-auto">
        <h2 className="mb-12 font-serif text-4xl font-bold text-center sm:text-5xl text-background">
          Experience Our Craft
        </h2>

        <div className="relative rounded-3xl overflow-hidden shadow-2xl h-96 sm:h-[500px] lg:h-[600px] bg-black">
          {!isPlaying ? (
            <>
              <video
                className="object-cover w-full h-full"
                poster="https://images.unsplash.com/photo-1495521821757-a1efb6729352?w=1200&h=800&fit=crop"
              >
                <source
                  src="https://www.w3schools.com/html/mov_bbb.mp4"
                  type="video/mp4"
                />
              </video>
              <div
                className="absolute inset-0 flex items-center justify-center transition-colors cursor-pointer bg-black/40 hover:bg-black/50"
                onClick={() => setIsPlaying(true)}
              >
                <button className="p-6 text-white transition-transform transform rounded-full bg-primary hover:bg-primary/90 hover:scale-110">
                  <Play size={40} fill="white" />
                </button>
              </div>
            </>
          ) : (
            <iframe
              className="w-full h-full"
              title="Hygee Coffee video"
              src="https://www.youtube.com/embed/N7isfDcCCJM?autoplay=1&mute=1&rel=0&modestbranding=1"
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
              referrerPolicy="strict-origin-when-cross-origin"
              allowFullScreen
            />
          )}
        </div>
      </div>
    </section>
  );
}
