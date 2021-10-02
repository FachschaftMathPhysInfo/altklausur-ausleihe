<template>
  <v-container>
    <div style="max-width: 960px; margin: 0px auto">
      <p class="text-h4">{{ $t('privacy.intro') }}</p>

      <p class="text-body1">
        {{ $t('privacy.text1') }}
      </p>

      <p class="text-h4">{{ $t('privacy.content') }}</p>
      <p class="text-h5">{{ $t('privacy.header1') }}</p>
<v-textarea
        readonly
        no-resize
        outlined
        rows="11"
        value="Fachschaft MathPhysInfo 
Im Neuenheimer Feld 205 (Mathematikon) 
69120 Heidelberg 
Deutschland 
fachschaft@mathphys.info 

Datenschutzbeauftragte
Verfasste Studierendenschaft der Universität Heidelberg
Albert-Ueberle-Str. 3-5 
69120 Heidelberg
datenschutz@stura.uni-heidelberg.de"
      >
      </v-textarea>
      <p class="text-h5">{{ $t('privacy.header2') }}</p>

      <p class="text-body1">
        {{ $t('privacy.text2') }}
        
      </p>

      <p class="text-body1">
        {{ $t('privacy.text3') }}
        
      </p>
      <p class="text-h5">{{ $t('privacy.header3') }}</p>

      <p class="text-body1">
        {{ $t('privacy.text4') }}
        
      </p>

      <p class="text-h6">{{ $t('privacy.subheader1') }}</p>
      <p class="text-body1">
        {{ $t('privacy.text5') }}
      </p>
      <p class="text-h6">{{ $t('privacy.subheader2') }}</p>
      <p class="text-body1">
        {{ $t('privacy.text6') }}
      
      <ul>
        <li>{{ $t('privacy.list1') }}</li>
        <li>{{ $t('privacy.list2') }}</li>
        <li>{{ $t('privacy.list3') }}</li>
        <li>{{ $t('privacy.list4') }}</li>
        <li>{{ $t('privacy.list5') }}</li>
      </ul>
      </p>
      <p class="text-h6">{{ $t('privacy.subheader3') }}</p>
      <p>{{ $t('privacy.text7') }}</p>
      
      <p class="text-h6">{{ $t('privacy.subheader4') }}</p>
      
      <p>{{ $t('privacy.text8') }}</p> 
      <p class="text-h6">{{ $t('privacy.subheader5') }}</p>
      <p> {{ $t('privacy.text9') }}
      
      <ul>
      <li>{{ $t('privacy.list6') }}</li> 
      <li>{{ $t('privacy.list7') }}</li></ul> </p>

      <p class="text-h5">{{ $t('privacy.header4') }}</p>
      <p> {{ $t('privacy.text10') }}</p>
    <p class="text-h5">{{ $t('privacy.header5') }}</p>
      <p> {{ $t('privacy.text11') }}</p>
      
      <v-textarea
        readonly
        no-resize
        outlined
        rows="6"
        value="Landesbeauftragte für Datenschutz und Informationsfreiheit Baden-Württemberg
Postfach 10 29 32, 70025 Stuttgart
Tel. 0711/615541-0; 
Fax 0711/615541-1
poststelle@lfdi.bwl.de 
https://www.baden-wuerttemberg.datenschutz.de"
      >
      </v-textarea>
      <p>{{ $t('privacy.text12') }}</p>
    </div>
  </v-container>
</template>

<script>
export default {
  name: "Privacy",

  data: () => ({})
};
</script>
