import { Button, Box } from "@chakra-ui/react";
import { FiSmartphone, FiUpload, FiPackage } from "react-icons/fi";
import { Link } from "wouter";
import { useInstallPack } from "../hooks/useInstallPack";
import { IsInstallingModal } from "./IsInstallingModal";
import { SyncMdMenu } from "./SyncMdMenu";

export const Header = () => {
  const { mutate: handleInstallPack, isLoading: isInstalling } =
    useInstallPack();
  return (
    <Box
      display="flex"
      position="fixed"
      backgroundColor="white"
      zIndex={2}
      top={0}
      right={0}
      left={0}
      p={2}
    >
      <IsInstallingModal isOpen={isInstalling} />
      <Box>
        <Link to="/device">
          <Button colorScheme="linkedin" leftIcon={<FiSmartphone />}>
            My Lunii
          </Button>
        </Link>
      </Box>
      <Box ml={2}>
        <SyncMdMenu />
      </Box>
      <Button
        variant="ghost"
        colorScheme="linkedin"
        leftIcon={<FiUpload />}
        onClick={handleInstallPack}
        ml={2}
      >
        Install pack
      </Button>
      <Link to="/create-pack">
        <Button
          variant="ghost"
          colorScheme="linkedin"
          leftIcon={<FiPackage />}
          ml={2}
        >
          Create pack
        </Button>
      </Link>
    </Box>
  );
};
