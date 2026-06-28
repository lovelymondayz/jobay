import { Canvas, useFrame } from "@react-three/fiber";
import { OrbitControls } from "@react-three/drei";
import { Color, AmbientLight, DirectionalLight, Clock } from "three";
import { FlipBook } from "quick_flipbook";
import { useEffect, useMemo } from "react";
import { Leva, useControls } from "leva";

function Book() {
  const clock = useMemo(() => new Clock(), []);
  const book = useMemo<FlipBook>(() => {
    const instance = new FlipBook({
      flipDuration: 0.7,
      yBetweenPages: 0.001,
      pageSubdivisions: 20,
    });
    instance.setPages([
      "pages/hygee-coffee-menu 2025_page-0001.webp",
      "pages/hygee-coffee-menu 2025_page-0002.webp",
      "pages/hygee-coffee-menu 2025_page-0003.webp",
      "pages/hygee-coffee-menu 2025_page-0004.webp",
      "pages/hygee-coffee-menu 2025_page-0005.webp",
      "pages/hygee-coffee-menu 2025_page-0006.webp",
      "pages/hygee-coffee-menu 2025_page-0007.webp",
      "pages/hygee-coffee-menu 2025_page-0008.webp",
      "pages/hygee-coffee-menu 2025_page-0009.webp",
      "pages/hygee-coffee-menu 2025_page-0010.webp",
      "pages/hygee-coffee-menu 2025_page-0011.webp",
      "pages/hygee-coffee-menu 2025_page-0012.webp",
      "pages/hygee-coffee-menu 2025_page-0013.webp",
      "pages/hygee-coffee-menu 2025_page-0014.webp",
      "pages/hygee-coffee-menu 2025_page-0015.webp",
      "pages/hygee-coffee-menu 2025_page-0016.webp",
      "pages/hygee-coffee-menu 2025_page-0017.webp",
      "pages/hygee-coffee-menu 2025_page-0018.webp",
      "pages/hygee-coffee-menu 2025_page-0019.webp",
      "pages/hygee-coffee-menu 2025_page-0020.webp",
      "pages/hygee-coffee-menu 2025_page-0021.webp",
      "pages/hygee-coffee-menu 2025_page-0022.webp",
    ]);
    return instance;
  }, []);

  // # for debug book size
  // const { debug, scaleX, scaleY, scaleZ } = useControls("Book Scale", {
  //   debug: false,
  //   scaleX: { value: 0.9, min: 0.1, max: 5, step: 0.01 },
  //   scaleY: { value: 2, min: 0.1, max: 5, step: 0.01 },
  //   scaleZ: { value: 1, min: 0.1, max: 5, step: 0.01 },
  // });

  useEffect(() => {
    function updateBookScale() {
      // if (debug) {
      //   book.scale.set(scaleX, scaleY, scaleZ);
      //   return;
      // }

      const isMobile = window.innerWidth < 768;
      if (isMobile) {
        book.scale.set(1.25, 2.0, 2.0);
      } else {
        book.scale.set(1.4, 3, 2.5);
      }
    }

    updateBookScale();
    window.addEventListener("resize", updateBookScale);
    return () => {
      window.removeEventListener("resize", updateBookScale);
    };
  }, [book]);
  // }, [book, debug, scaleX, scaleY, scaleZ]);

  useFrame(() => {
    const delta = clock.getDelta();
    book.animate(delta);
  });

  useEffect(() => {
    const nextPage = () => book.nextPage();
    const previousPage = () => book.previousPage();
    const nextBtn = document.querySelector<HTMLButtonElement>(".next");
    const prevBtn = document.querySelector<HTMLButtonElement>(".prev");
    nextBtn?.addEventListener("click", nextPage);
    prevBtn?.addEventListener("click", previousPage);
    return () => {
      nextBtn?.removeEventListener("click", nextPage);
      prevBtn?.removeEventListener("click", previousPage);
      book.dispose();
    };
  }, [book]);

  return <primitive object={book} />;
}

export default function FlipBookScene() {
  return (
    <>
      <Leva collapsed={false} />
      <Canvas
        className="w-full h-full"
        camera={{ position: [0, 1.2, 4.5], fov: 45, near: 0.1, far: 1000 }}
        gl={{ antialias: true, alpha: true }}
        onCreated={({ scene }) => {
          scene.background = new Color(0x0f172a);
        }}
      >
        <primitive object={new AmbientLight(0xffffff, 2)} />
        <primitive
          object={new DirectionalLight(0xffffff, 3)}
          position={[5, 5, 5]}
        />
        <OrbitControls
          enableDamping
          enablePan={false}
          minDistance={2}
          maxDistance={8}
          maxPolarAngle={Math.PI / 2}
        />
        <Book />
      </Canvas>
    </>
  );
}
