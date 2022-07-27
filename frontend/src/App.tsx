import { Box, Button, Center, Icon, Text, Tooltip } from "@chakra-ui/react";
import { useQuery } from "@tanstack/react-query";
import { BiPlug } from "react-icons/bi";
import {
  FiPackage,
  FiRefreshCcw,
  FiSmartphone,
  FiUpload,
} from "react-icons/fi";
import { Link, Route } from "wouter";
import { ListPacks } from "../wailsjs/go/main/App";
import { DetailsModal } from "./components/DetailsModal";
import { DeviceModal } from "./components/DeviceModal";
import { Header } from "./components/Header";
import { IsInstallingModal } from "./components/IsInstallingModal";
import { NewPackModal } from "./components/NewPackModal";
import { PackList } from "./components/PackList";
import { SyncMdMenu } from "./components/SyncMdMenu";
import { useDeviceQuery } from "./hooks/useApi";
import { useInstallPack } from "./hooks/useInstallPack";

function App() {
  const {
    data: device,
    refetch: refetchDevice,
    isLoading: isLoadingDevice,
  } = useDeviceQuery();

  if (isLoadingDevice) return null;

  return (
    <Box id="App" p={3}>
      <Route path="/create-pack">
        <NewPackModal />
      </Route>
      <Route path="/device">
        <DeviceModal />
      </Route>
      <Route path="/pack/:id">
        <DetailsModal />
      </Route>

      {!device && (
        <>
          <Box display="flex">
            <Tooltip label="Try refreshing after the device has been mounted by the system">
              <Button
                colorScheme="orange"
                rightIcon={<FiRefreshCcw />}
                onClick={() => refetchDevice()}
              >
                Refresh
              </Button>
            </Tooltip>
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
          <Center opacity={0.3} fontSize={20} h={500} flexDirection="column">
            <Icon as={BiPlug}></Icon>

            <Text>No device connected</Text>
          </Center>
        </>
      )}

      {device && (
        <Box pt="56px">
          <Header />
          <PackList />
        </Box>
      )}
    </Box>
  );
}

export default App;
