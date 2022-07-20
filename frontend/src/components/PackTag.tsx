import { Tooltip, Tag } from "@chakra-ui/react";
import { FC } from "react";
import { lunii } from "../../wailsjs/go/models";

export const PackTag: FC<{ metadata: lunii.Metadata }> = ({ metadata }) => {
  if (!metadata.isUnknown && metadata.isOfficialPack)
    return (
      <Tooltip label="This pack was found in lunii's official database">
        <Tag colorScheme="teal">Official</Tag>
      </Tooltip>
    );
  if (!metadata.isUnknown && !metadata.isOfficialPack)
    return (
      <Tooltip label="This pack seems to be a custom made pack">
        <Tag colorScheme="blue">Custom</Tag>
      </Tooltip>
    );
  if (metadata.isUnknown)
    return (
      <Tooltip label="This pack might have been installed with STUdio. In a future version, this tool will be able to get metadatas for those packs too.">
        <Tag colorScheme="orange">No Metadata</Tag>
      </Tooltip>
    );
  return null;
};
