import React from 'react';
import { Progress } from "@nextui-org/react";
import { keyframes } from "@emotion/react";
import styled from "@emotion/styled";
import Logo from "@/app/assets/logo_full.png"

const waveAnimation = keyframes`
  0% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
  100% { transform: translateY(0); }
`;

const StyledLogo = styled.img`
  animation: ${waveAnimation} 2s ease-in-out infinite;
`;

const LoadingScreen = () => {

  return (
    <div className="fixed inset-0 flex flex-col items-center justify-center bg-white dark:bg-black z-50">
      <StyledLogo
        src={Logo.src}
        alt=""
        className="h-16 md:h-24 mb-8"
      />
      <div className="w-64 px-4">
        <Progress
          size="sm"
          isIndeterminate
          aria-label="Loading..."
          className="max-w-md"
          color="primary"
        />
      </div>
    </div>
  );
};

export default LoadingScreen;