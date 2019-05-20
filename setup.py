from distutils.core import setup

setup(
    name='fishbeats',
    version='0.0.1',
    packages=['fishbeats'],
    url='github.com/AlexisBRENON/fishbeats',
    license='',
    author='Alexis BRENON',
    author_email='brenon.alexis+fishbeats@gmail.com',
    description='Generate music from moving fishes in an aquarium',
    package_dat={
        'configuration': ['share/']
    },
    install_requires=[
        "pyFluidSynth == 1.*",
        "Autologging == 1.*",
        "protobuf == 3.*",
    ]
)
