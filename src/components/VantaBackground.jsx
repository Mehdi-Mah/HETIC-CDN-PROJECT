import React, { useEffect, useRef, useState } from "react";
import * as THREE from "three";
import CLOUDS from "vanta/dist/vanta.clouds.min";

const VantaBackground = () => {
	const [vantaEffect, setVantaEffect] = useState(null);
	const vantaRef = useRef(null);

	useEffect(() => {
		if (!vantaEffect) {
			setVantaEffect(
				CLOUDS({
					el: vantaRef.current,
					THREE,
					mouseControls: true,
					touchControls: true,
					gyroControls: false,
					minHeight: 200.00,
					minWidth: 200.00,
					skyColor: 0x67b8d9,
					cloudColor: 0xadc2de,
					cloudShadowColor: 0x1c3b57,
					sunColor: 0xfa9a23,
					sunGlareColor: 0xf76c3d,
					sunlightColor: 0xfc9834,
					mouseEase: 1,
					speed: 0.90
				})
			);
		}

		return () => {
			if (vantaEffect) vantaEffect.destroy();
		};
	}, [vantaEffect]);

	return <div ref={vantaRef} className="vanta-container"></div>;
};

export default VantaBackground;
