import { Tooltip, Tag } from "@chakra-ui/react";
import { FC } from "react";
import { lunii } from "../../wailsjs/go/models";

export const PackTag: FC<{ metadata: lunii.Metadata }> = ({ metadata }) => {
  if (metadata.packType === "lunii")
    return (
      <Tooltip label="This pack was found in lunii's official database">
        <Tag colorScheme="teal">Lunii</Tag>
      </Tooltip>
    );
  if (metadata.packType === "custom")
    return (
      <Tooltip label="This pack is a custom made pack">
        <Tag colorScheme="blue">Custom</Tag>
      </Tooltip>
    );
  if (metadata.packType === "undefined")
    return (
      <Tooltip label="This pack might have been installed from STUdio from another computer">
        <Tag colorScheme="orange">No Metadata</Tag>
      </Tooltip>
    );
  return null;
};
